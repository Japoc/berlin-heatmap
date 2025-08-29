package gql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	url    string
	client *http.Client
}

func New(url string) *Client {
	return &Client{
		url:    url,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type Stop struct {
	ID          string
	Lat         float64
	Lon         float64
	VehicleMode string
}

func (c *Client) Stops(ctx context.Context, box BBox) ([]Stop, error) {
	query := `{
		stops {
			id
			lat
			lon
			vehicleMode
		}
	}`
	var resp struct {
		Data struct {
			Stops []struct {
				ID          string  `json:"id"`
				Lat         float64 `json:"lat"`
				Lon         float64 `json:"lon"`
				VehicleMode string  `json:"vehicleMode"`
			} `json:"stops"`
		} `json:"data"`
	}
	if err := c.do(ctx, query, nil, &resp); err != nil {
		return nil, err
	}
	stops := make([]Stop, len(resp.Data.Stops))
	count := 0
	for _, s := range resp.Data.Stops {
		if isInBox(Stop(s), box) && s.VehicleMode == "RAIL" {
			stops[count] = Stop{s.ID, s.Lat, s.Lon, s.VehicleMode}
			count++
		}
	}
	log.Printf(`found %v stops`, count)
	return stops[:count], nil
}

type Route struct {
	Mode      string
	ShortName string
	Color     string
	Points    string
}

func (c *Client) Routes(ctx context.Context, routeNames []string, mode string) ([]Route, error) {
	routes := make([]Route, 0)
	for _, routeName := range routeNames {
		query := fmt.Sprintf(`{
		  routes(name: "%v") {
			mode
			shortName
			type
			textColor
			color
			patterns {
			  stops {
				name
			  }
			  patternGeometry {
				points
			  }
			}
		  }
		}`, routeName)
		var resp struct {
			Data struct {
				Routes []struct {
					Mode      string `json:"mode"`
					ShortName string `json:"shortName"`
					Type      int    `json:"type"`
					TextColor string `json:"textColor"`
					Color     string `json:"color"`
					Patterns  []struct {
						Stops []struct {
							Name string `json:"name"`
						}
						PatternGeometry struct {
							Points string `json:"points"`
						}
					} `json:"patterns"`
				} `json:"routes"`
			} `json:"data"`
		}

		if err := c.do(ctx, query, nil, &resp); err != nil {
			return nil, err
		}
		for _, r := range resp.Data.Routes {
			if r.Mode == mode {
				mostStops := 0
				points := ""
				for _, pattern := range r.Patterns {
					if len(pattern.Stops) > mostStops {
						mostStops = len(pattern.Stops)
						points = pattern.PatternGeometry.Points
					}
				}
				routes = append(routes, Route{
					Mode:      r.Mode,
					ShortName: r.ShortName,
					Color:     r.Color,
					Points:    points,
				})
			}
		}
		log.Printf(`found %v routes`, len(routes))
	}

	return routes, nil
}

func (c *Client) TravelTime(ctx context.Context, from, to Stop, when time.Time) (int, error) {
	query := `query($fromLat: CoordinateValue!, $fromLon: CoordinateValue!, $toLat: CoordinateValue!, $toLon: CoordinateValue!, $dt: OffsetDateTime!) {
	  planConnection(
	    origin: {
        location: {coordinate: {
           latitude: $fromLat, longitude: $fromLon
        }}
      },
	    destination: {
        location: {coordinate: {
           latitude: $toLat, longitude: $toLon
        }}
       
      },
	    dateTime: { earliestDeparture: $dt},
	    modes: {direct: WALK, transit: {transit: [{mode:BUS}, {mode:RAIL}, {mode:TRAM}]}}
	  ) {
	    edges {
	      node {
          duration
        }
	    }
	  }
	}`
	vars := map[string]interface{}{
		"fromLat": from.Lat, "fromLon": from.Lon,
		"toLat": to.Lat, "toLon": to.Lon,
		"dt": when.Format(time.RFC3339),
	}
	var resp struct {
		Data struct {
			PlanConnection struct {
				Edges []struct {
					Node struct {
						Duration int `json:"duration"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"planConnection"`
		} `json:"data"`
	}
	if err := c.do(ctx, query, vars, &resp); err != nil {
		return -1, err
	}
	if len(resp.Data.PlanConnection.Edges) == 0 {
		return -1, fmt.Errorf("no itinerary found")
	}
	// choose shortest duration (seconds → minutes)
	min := resp.Data.PlanConnection.Edges[0].Node.Duration
	for _, it := range resp.Data.PlanConnection.Edges {
		if it.Node.Duration < min {
			min = it.Node.Duration
		}
	}
	return min / 60, nil
}

func (c *Client) do(ctx context.Context, query string, vars map[string]interface{}, out interface{}) error {
	body, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": vars,
	})
	req, _ := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status: %d %s", resp.StatusCode, string(b))
	}
	return json.Unmarshal(b, out)
}

func isInBox(from Stop, box BBox) bool {
	lat1, lat2, long1, long2 := box.MinLat, box.MaxLat, box.MinLon, box.MaxLon
	if lat1 < lat2 {
		lat1, lat2 = lat2, lat1
	}
	if long1 < long2 {
		long1, long2 = long2, long1
	}
	return (from.Lon >= long2 && from.Lon <= long1) && (from.Lat >= lat2 && from.Lat <= lat1)
}

type BBox struct {
	MinLon, MinLat, MaxLon, MaxLat float64
}
