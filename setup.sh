curl -X PUT -H 'Content-Type: application/json' --data '{"steam":{"SteamPath":"D:/Games/Steam"},"mesh":{"peers":[{"name": "sentinel-prime"}]}}' http://bigred:9771/app/config
curl -X PUT -H 'Content-Type: application/json' --data '{"steam":{"SteamPath":"/home/lucas/.steam/steam"},"mesh":{"peers":[{"name": "sentinel-prime"}]}}' http://192.168.14.18:9771/app/config
curl -X PUT -H 'Content-Type: application/json' --data '{"mesh":{"peers":[{"name": "bigred"},{"name": "192.168.14.18"}]}}' http://localhost:9771/app/config
curl -X POST http://localhost:9771/mesh/copy/bigred/270370
