# Messages Websocket

## Connexion au Websocket

Pour se connecter au websocket permettant de commencer une session de navigation, il faut faire une requête sur l'endpoint suivant :

`http://addresse_supmap/navigation/ws?session_id=XXXXXX`

Le paramètre `session_id` contient un UUID généré par le client, il permet d'identifier la session de navigation active.

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

## Emits par le serveur

### Incident

Type : `incident`

Ce message est envoyé par le serveur dès qu'une action liée aux incidents ("create", "deleted", "certified") est réalisée par le microservice **supmap-incidents**.  
Il informe tous les clients connectés d'un changement concernant un incident sur leur trajet respectif.

Un incident peut représenter, par exemple, un embouteillage, un accident, ou tout autre événement susceptible d'impacter la circulation.

Le champ `action` précise le type d'opération réalisée sur l'incident :
* `"create"` : un nouvel incident a été détecté et ajouté.
* `"certified"` : l’incident a été confirmé (nombre requis d'intéractions atteint).
* `"deleted"` : l’incident n’est plus d’actualité.

Le champ `incident` contient les informations détaillées sur l’incident concerné.

Exemple :

```json
{
  "type": "incident",
  "data": {
    "incident": {
      "id": 26,
      "user_id": 2,
      "type": {
        "id": 3,
        "name": "Embouteillage",
        "description": "Circulation fortement ralentie ou à l’arrêt.",
        "need_recalculation": true
      },
      "lat": 49.19477822,
      "lon": -0.3964915,
      "created_at": "2025-05-09T14:57:36.96141Z",
      "updated_at": "2025-05-09T14:57:36.96141Z"
    },
    "action": "create"
  }
}
```

#### Détail des champs de `data` :

- `incident` : Objet décrivant l’incident.
  - `id` : Identifiant unique de l’incident.
  - `user_id` : Identifiant de l’utilisateur ayant signalé l'incident.
  - `type` : Objet précisant la nature de l’incident.
    - `id` : Identifiant du type d’incident.
    - `name` : Nom du type d’incident (ex : "Embouteillage").
    - `description` : Description détaillée du type d’incident.
    - `need_recalculation` : Booléen indiquant si la présence de ce type d’incident nécessite le recalcul de l’itinéraire.
  - `lat` : Latitude de l’incident.
  - `lon` : Longitude de l’incident.
  - `created_at` : Date de création de l’incident (format ISO8601, UTC).
  - `updated_at` : Date de la dernière mise à jour de l’incident.
  - `deleted_at` _(optionnel)_ : Date de suppression de l’incident (présent uniquement si l’incident est supprimé).

- `action` : Type d’action liée à l’incident. Peut être `"create"`, `"certified"` ou `"deleted"`.

---

_Note : Les incidents sont transmis en temps réel. Le client doit adapter son comportement selon le type d’action reçue (affichage, recalcul d’itinéraire, suppression, etc.)._
