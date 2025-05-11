# supmap-navigation

## 1. Introduction

### 1.1. Rôle du microservice

**supmap-navigation** est le microservice dédié à la gestion de la navigation en temps réel pour les utilisateurs de l’application Supmap. Il établit et maintient des connexions WebSocket avec les clients mobiles afin de :
- Suivre en direct la position de chaque utilisateur pendant leur trajet.
- Diffuser instantanément les nouveaux incidents signalés sur leur itinéraire.
- Gérer les recalculs de route à la volée en cas d’événement perturbateur (ex : accident, embouteillage).

### 1.2. Principales responsabilités

- **Connexion WebSocket et gestion de session :**  
  Chaque client ouvre une connexion WebSocket identifiée par un `session_id` unique (UUID). Le serveur conserve en cache les informations de navigation et les positions des clients grâce à Redis.

- **Suivi de position :**  
  Les clients envoient régulièrement leur position. Le service met à jour le cache et peut ainsi déterminer à tout moment l’avancement de l’utilisateur sur son trajet.

- **Diffusion d’incidents en temps réel :**  
  Lorsqu’un nouvel incident est détecté ou modifié (via le microservice supmap-incidents), supmap-navigation est notifié via un canal Pub/Sub Redis. Il transmet alors en temps réel l’incident aux clients concernés, c’est-à-dire ceux dont l’itinéraire croise la zone de l’incident.

- **Recalcul dynamique des itinéraires :**  
  Si un incident nécessite le recalcul de la route (incident bloquant et certifié…), le service interroge supmap-gis pour obtenir un nouvel itinéraire. Ce nouvel itinéraire est ensuite envoyé au(x) client(s) via la connexion WebSocket, assurant une navigation optimisée en permanence.

### 1.3. Technologies et dépendances externes

- **Go** : langage principal du microservice.
- **WebSocket** : communication temps réel bidirectionnelle avec les clients.
- **Redis** : 
  - Stockage temporaire (cache) des sessions, routes et positions clients.
  - Mécanisme Pub/Sub pour recevoir en direct les incidents depuis supmap-incidents.
- **supmap-gis** : microservice utilisé pour le recalcul d’itinéraires en cas d’incident bloquant.
- **supmap-incidents** : source des incidents signalés sur le réseau via Redis Pub/Sub.
- **GitHub Actions** : CI pour build/push l’image Docker sur le registre GHCR du repo.

---

