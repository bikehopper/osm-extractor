# OSM Extractor

Makes OpenStreetMap cutouts from the larger USA data set.

## Dependencies

* [docker](https://www.docker.com/)
* [node.js](https://nodejs.org/en)

1. If you do not plan to run this locally, skip this step. Find and download the smallest `.osm.pbf` file that includes the entire area you want to extract from [Geofabrik's download page](https://download.geofabrik.de/north-america/us.html) (This example uses the Northern California file). This may take a while, so you can start this download, then move on to the next steps. 
```
curl 'https://download.geofabrik.de/north-america/us/california/norcal-latest.osm.pbf' -o ./volumes/input/latest.osm.pbf
```
2. Decide which counties you want to include in your OSM cutout. Find each county in [counties-usa.geojson](counties-usa.geojson), confirm spelling, capitalization, etc...
3. Add all the counties you want in your extract as a comma seperated list in the following command, e.g. `Alameda,Contra Costa,Marin,Napa,San Mateo,Santa Clara,Solano,Sonoma,San Francisco`. Set the name of the output file. Then run it. This should create a geojson file in the `./polygons` directory. 
```
npx --yes mapshaper \
    -i counties-usa.geojson \
    -filter '"{list-of-comma-separated-county-names}".indexOf(NAME) > -1' \
    -dissolve2 \
    -o ./polygons/{my-region-name}.geojson geojson-type=Feature
```
4. (optional) Create a convex hull of the counties. This is useful when the group of counties don't make a solid shape. Use the geojson file created in the previous step to make a new geojson file. 
```
npx --yes turf-cli convex polygons/san-francisco-bay-area.geojson > polygons/san-francisco-bay-area-convex.geojson
```
5. Open [config.json](config.json) and add an entry for the extract you'd like to create. Copy an existing entry and replace the values accordingly.
6. (test) Build the docker container by running:
```
docker build . -t osm-extractor:testing-new-extract
```
7. (test) Generate the extract by running the docker container. This should create a new .pbf file in `./volumes/output`.
```
docker run \
-v ./volumes/output:/mnt/output \
-v ./volumes/input:/mnt/input \
osm-extractor:testing-new-extract osmium extract \
	-d /mnt/output \
	-c /app/config.json \
	/mnt/input/latest.osm.pbf
```
