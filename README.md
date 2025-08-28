http://localhost:8088/heatmap?lat=52.5208&lon=13.4095&format=json
http://localhost:8088/heatmap?lat=52.5208&lon=13.4095&format=png

web:
http://localhost:5173/

run otp:
````
docker run -it --rm -p 8080:8080 -m=15g -v ./berlin:/var/opentripplanner docker.io/opentripplanner/opentripplanner:latest --load --serve
````
otp: 
http://localhost:8080/
https://docs.opentripplanner.org/en/latest/apis/GraphQL-Tutorial/

graphql web client:
http://localhost:8080/graphiql?query=query%20%7B%0A%20%20trip1%3A%20plan(%0A%20%20%20%20from%3A%20%7Blat%3A%2052.5200%2C%20lon%3A%2013.4050%7D%2C%20%23%20Berlin%20Mitte%0A%20%20%20%20to%3A%20%20%20%7Blat%3A%2052.5300%2C%20lon%3A%2013.4000%7D%2C%20%23%20Nearby%20location%0A%20%20%20%20date%3A%20%222025-08-22T12%3A00%3A00%2B02%3A00%22%2C%0A%20%20%20%20transportModes%3A%20%5B%7Bmode%3A%20BUS%7D%2C%20%7Bmode%3A%20RAIL%7D%5D%0A%20%20)%20%7B%0A%20%20%20%20itineraries%20%7B%0A%20%20%20%20%20%20duration%0A%20%20%20%20%20%20walkTime%0A%20%20%20%20%20%20waitingTime%0A%20%20%20%20%7D%0A%20%20%7D%0A%0A%20%20trip2%3A%20plan(%0A%20%20%20%20from%3A%20%7Blat%3A%2052.5150%2C%20lon%3A%2013.3900%7D%2C%20%23%20Potsdamer%20Platz%0A%20%20%20%20to%3A%20%20%20%7Blat%3A%2052.5400%2C%20lon%3A%2013.4100%7D%2C%20%23%20Prenzlauer%20Berg%0A%20%20%20%20date%3A%20%222025-08-22T12%3A00%3A00%2B02%3A00%22%2C%0A%20%20%20%20transportModes%3A%20%5B%7Bmode%3A%20BUS%7D%2C%20%7Bmode%3A%20RAIL%7D%5D%0A%20%20)%20%7B%0A%20%20%20%20itineraries%20%7B%0A%20%20%20%20%20%20duration%0A%20%20%20%20%20%20walkTime%0A%20%20%20%20%20%20waitingTime%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D%0A

query {
trip1: plan(
from: {lat: 52.5200, lon: 13.4050}, # Berlin Mitte
to:   {lat: 52.5300, lon: 13.4000}, # Nearby location
date: "2025-08-22T12:00:00+02:00",
transportModes: [{mode: BUS}, {mode: RAIL}]
) {
itineraries {
duration
walkTime
waitingTime
}
}

trip2: plan(
from: {lat: 52.5150, lon: 13.3900}, # Potsdamer Platz
to:   {lat: 52.5400, lon: 13.4100}, # Prenzlauer Berg
date: "2025-08-22T12:00:00+02:00",
transportModes: [{mode: BUS}, {mode: RAIL}]
) {
itineraries {
duration
walkTime
waitingTime
}
}
}
