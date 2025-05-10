# Messages Websocket

## Connexion au Websocket

Pour se connecter au websocket permettant de commencer une session de navigation, il faut faire une requête sur l'endpoint suivant :

`http://addresse_supmap/navigation/ws?session_id=XXXXXX`

Le paramètre `session_id` contient un UUID généré par le client et stocké dans le local storage, il permet d'identifier la session de navigation active.

## Structure générale

Les messages échangés entre le client et le serveur sont au format JSON et respectent la structure suivante :

```json
{
  "type": TYPE_MESSAGE,
  "data": {...}
}
```

Les types de message sont les suivants :
* Emits par le serveur :
  * "route"
  * "incident"
* Emits par le client :
  * "init"
  * "position"

Le champ `data` est un objet qui dépend du type de message.

## Emits par le client

### Initialisation

Type : `init`

Ce message est envoyé par le client dès qu'il se connecte au websocket. Il permet au serveur de connaître la route du client et sa position à l'instant de la connexion.

Exemple :

```json
{
    "type": "init",
    "data": {
        "session_id": "e38d5757-5359-44b3-ab6e-8c619e3daba0",
        "last_position": {
            "lat": 49.171669,
            "lon": -0.582579,
            "timestamp": "2025-05-06T20:52:30Z"
        },
        "route": {
            "polyline": [
                {
                    "latitude": 49.171669,
                    "longitude": -0.582579
                },
                .........,
                {
                    "latitude": 49.201345,
                    "longitude": -0.392996
                }
            ],
            "locations": [
                {
                    "lat": 49.17167279051877,
                    "lon": -0.5825858234777268
                },
                {
                    "lat": 49.20135359834111,
                    "lon": -0.3930605474075204
                }
            ]
        },
        "updated_at": "2025-05-06T20:52:30Z"
    }
}
```

Le champ `session_id` doit contenir le même UUID utilisé à la connexion initiale (en HTTP).

Le champ `locations` correspond aux points d'arrêts (départ, arrivée et points intermédiaires s'il y en a) de l'itinéraire.

### Position

Type : `position`

Ce message est envoyé à intervalles réguliers par le client (cinq secondes). Il contient la position actuelle du client accompagné d'un timestamp. Ainsi, le serveur est toujours au courant de la position actuelle du client.

Exemple : 

```json
{
    "type": "position",
    "data": {
        "lat": 49.1943057668118,
        "lon": -0.44595408906894096,
        "timestamp": "2025-05-07T10:07:00Z"
    }
}
```
