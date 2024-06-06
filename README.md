# OSM Extractor

Makes OpenStreetMap cutouts from the larger USA data set.

## Dependencies

* [docker](https://www.docker.com/)
* [node.js](https://nodejs.org/en)

## Usage

Add a geojson file that contains a Feature into the `./polygons` directory, add a corresponding entry in `./config.json`, then build and run the docker container.

Follow the steps below to use the tool to create your own cutouts. The San Francisco Bay Area consists of nine counties; they are used in the following example.

1. If you do not plan to run this locally, skip this step. Find and download the smallest `.osm.pbf` file that includes the entire area you want to extract from [Geofabrik's download page](https://download.geofabrik.de/north-america/us.html) (This example uses the Northern California file). This may take a while, so you can start this download, then move on to the next steps. 
    ```
    curl 'https://download.geofabrik.de/north-america/us/california/norcal-latest.osm.pbf' -o ./volumes/input/latest.osm.pbf
    ```
2. Do one of the two following options:
    1. Decide what counties you want to include in your OSM cutout. Find each county in [counties-usa.geojson](counties-usa.geojson), confirm spelling, capitalization, etc... Add all of the counties you want in your extract as a comma seperated list in the follow command, replacing `Alameda,Contra Costa,Marin,Napa,San Mateo,Santa Clara,Solano,Sonoma,San Francisco`. Set the name of the output file (replacing `bay-area-geometry.geojson` with something that makes sense). Then run it. This should create a geojson file in the `./polygons` directory. 
        ```
        npx --yes mapshaper -i counties-usa.geojson -filter '"Alameda,Contra Costa,Marin,Napa,San Mateo,Santa Clara,Solano,Sonoma,San Francisco".indexOf(NAME) > -1' -dissolve2 -o ./polygons/san-francisco-bay-area.geojson geojson-type=Feature
        ```

    2. Use a website like [geojson.io](https://geojson.io/) to make a polygon, then same the json to the `./polygons` directory. 

3. [optional] Create a convex hull of the counties. This is useful when the group of counties don't make a solid shape. Use the geojson file created in the previous step to make a new geojson file. 
    ```
    npx --yes turf-cli convex polygons/san-francisco-bay-area.geojson > polygons/san-francisco-bay-area-convex.geojson
    ```

4. Open [config.json](config.json) and add an entry for the extract you'd like to create. Copy an existing entry and replace the values accordingly.
5. [testing] Build the docker container by running:
    ```
    docker build . -t osm-extractor:testing-new-extract
    ```
6. [testing] Generate the extract by running the docker container. This should create a new .pbf file in `./volumes/output`.
    ```
    docker run -v ./volumes/output:/mnt/output -v ./volumes/input:/mnt/input osm-extractor:testing-new-extract osmium extract -d /mnt/output -c /app/config.json /mnt/input/latest.osm.pbf
    ```
