


<h3 align="center">Berlin Public Transport Heatmap</h3>

  <p align="center">
    A tool to visualize travel times around berlin using public transport
  </p>
<p align="center">
    A hosted version can be found here: <a href="https://japoc.github.io/berlin-heatmap/">https://japoc.github.io/berlin-heatmap/</a>
  </p>


<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#notes">Notes</a></li>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#precompute">Precompute</a></li>
        <li><a href="#serve">Serve</a></li>
      </ul>
    </li>
  </ol>
</details>


### About The Project

![Heatmap Screen Shot](./static/screenshot.png)

This project is a tool to visualize travel times around berlin using public transport.

### Notes

The project currently only includes one dataset, based on the transit data on 22.08.2025 at 12:00.

The project currently only includes the about 700 train stations of Berlin, so no tram or bus are being considered.

### Built With

+ frontend:  
  + vite
  + vue
  + typescript
  + leaflet for the map and overlaying it with a generated heatmap (png)
+ backend: golang
+ data: 
  + osm data for berlin: [osm](https://download.geofabrik.de/europe/germany/berlin-latest.osm.pbf)
  + gtfs data for berlin [vbb](https://vbb.de/vbbgtfs)
  + tool to compute distances between stops: [Open Trip Planner](https://docs.opentripplanner.org/en/latest/)

  
## Getting Started

```
 # create directory for data and config
 mkdir berlin
 # download OSM
 curl -L https://download.geofabrik.de/europe/germany/berlin-latest.osm.pbf -o berlin/osm.pbf  
 # download GTFS
 curl -L https://vbb.de/vbbgtfs -o berlin/vbb-gtfs.zip
 # build graph and save it onto the host system via the volume
 docker run --rm \
     -e JAVA_TOOL_OPTIONS='-Xmx8g' \
     -v "$(pwd)/berlin:/var/opentripplanner" \
     docker.io/opentripplanner/opentripplanner:latest --build --save
 # load and serve graph
 docker run -it --rm -p 8080:8080 \
     -e JAVA_TOOL_OPTIONS='-Xmx8g' \
     -v "$(pwd)/berlin:/var/opentripplanner" \
     docker.io/opentripplanner/opentripplanner:latest --load --serve
```

This downloads both the OSM and GTFS and builds a graph that is needed to run an OTP instance.
If everything was successful the OTP Debug UI can be accessed on http://localhost:8080

The OTP instance also comes with a GTFS GraphQL WebClient, which can be accessed at http://localhost:8080/graphiql 

This graphQL Client is necessary to precompute datasets to generate and display heatmaps at runtime.

### Precompute

To precompute use the `scripts/precompute.sh` script.
It supports the following flags:
* `graphql`: url to the GraphQL, in most cases should be http://localhost:8080/otp/gtfs/v1
* `bucket`: name of the dataset (used for output)
* `bbox`: Bounding box of the area you are calculating for as latitude and logitude, example for berlin: `"13.0884,52.3383,13.7612,52.6755"`
* `grid`: size of the grid used to calculate the heat map in meters (default = 400)
* `k`: amount of the nearest stops from a given coordinate that are used to find the quickest travel time
* `out`: path where to output the precomputed dataset
* `datetime`: datetime for which the dataset will be created, default: `2025-08-22T12:00:00Z`

precompute will create:
+ `stops.csv`: a list of all stops within the bounding box, with latitude and longitude
+ `matrix_{bucket}`: distance from every stop to another, as binary
+ `grid_links.json`: a list of the k nearest stops for every grid, and their walking distance

Precomputation may take some time, as it is required to make (O)² GraphQL Queries (O = Amount of Stops to be considered)
+ Only Trainstops (around 800) would mean 800*800 = 640.000 queries
+ Adding bus stops on top (around 6500) would mean = 53.290.000 queries

### Serve
To serve the frontend locally run `npm run dev` in `/web`

To serve the backend locally use `serve.sh` `/scripts`
It supports the following flags:
* `port`: port where the backend will be served, default is 8088
* `matrix`: path to the precomputed matrix 
* `gridFile`: path to the precomputed `grid_links.json`
* `stops`: path to the precomputed `stops.json`
