# Praxisprojekt

## Einleitung

Dieser Repo beinhaltet den source code, der im Rahmen der Praxisphase zustande gekommen ist. Die Praxisphase beinhaltete
die Evaluierung der [Matrix Spezifikation](https://spec.matrix.org/latest/) und die Realisierung eines in diesem Rahmen
festgelegten PoCs(Proof of Concept). Für weitere Anforderungsanalysen siehe den
[mural board](https://app.mural.co/t/sprinteins1549/m/sprinteins1549/1636105317890/65870bfda3ff79f8b53c2345e2c8d79a2ef7f938?sender=15c9e8be-0d6a-418c-80b6-0b89ec5eb49d).
Es ist noch zu erwähnen, dass im Rahmen des Praxisprojekts die *client-server API* Anwendung fand. 

### Proof of Concept

Als PoC wurden folgende Themen priorisiert:

- Standard-Implementierung von Frontends und Backends
  - [Element](https://github.com/vector-im/element-web) als FE
  - [Dendrite](https://github.com/matrix-org/dendrite) als erstes Options für BE
  - [Synapse](https://github.com/matrix-org/synapse) als zweite Option für BE
- Deployment der Standard-Implementierungen in k8s via helm charts
- SSO
- Automatisiertes Raum-Management
- Vergabe der notwendigen Rechte an User in einem Raum

Anschließend kann sich daraus das folgende Szenario abgeleitet werden:

___Der User kann sich via SSO (oder testweise durch einen externen OpenID-Connect Anbieter) durch FE einloggen und beim
BE als Matrix-User hinterlegt werden. Des Weiteren wird dieser User automatisiert in einem zuvor automatisiert erstellten
Raum mit bestimmten Rechten zugeordnet.___

## Struktur

In dem top directory befinden sich pro services einen entsprechenden Ordner mit demselben Namen wie der Service-Name.
Der Ordner `docker` dient als erste Testimplementierung mit _docker_ und _docker compose_ und kann ignoriert werden. Der
Ordner `middleware` stellt einen service dar, der als _admin bot_ zur Raumerstellung und Rechteverteilung angesehen
werden kann.

Die anderen Ordner beinhalten die jeweiligen _helm chart_ Implementierungen.

### Standard-BEs

Da bei dem Matrix-BE zwei verschiedene Server-Implementierung gibt, galt diese miteinander zu vergleich und einzeln zu
deployen. Aus diesem Grund befinden zwei Ordner mit den zwei BE-Namen: `dendrite` und `synapse`. Nachfolgend werden die
beiden Implementierungen miteinander verglichen und das Resultat tabellarisch aufgelistet (Stand 23.01.22):
 
| Features       | Dendrite                                   | Synapse      |
|:---------------|:-------------------------------------------|:-------------|
| version        | v0.5.1                                     | v1.49.2      |
| adoption       | beta, early access                         | stable       |
| language       | go                                         | python       |
| mode           | monolith, polylith                         | monolith     |
| db             | psql, sqlite                               | psql, sqlite |
| event handling | kafka, nafka (custom kafka implementation) | internal     |
| user limits    | 10s - 100s                                 | -            |
| scalability    | only with polylith mode and kafka          | internal     |
| restriction    | SSO not implemented yet                    | -            |
| FE connection  | successful                                 | successful   |

Aus dem obigen Vergleich wird deutlich, dass der `dendrite` Server nicht zur Erfüllung des PoCs eingesetzt werden kann,
da er den SSO noch nicht implementiert <sup> [1](https://github.com/matrix-org/dendrite/pull/2014) </sup>.
Für die Evaluierung wurde dennoch der `dendrite` server in _k8s_ deployed und mit FE verbunden. Für die weiteren Schritte
wurde der `synapse` server näher betrachtet.

## Anforderungen

Diese Anforderungen __müssen__ vorhanden sein, dennoch werden sie mithilfe des `init` Script installiert werden.
Ansonsten müssen diese bereits auf dem Zielbetriebssystem bereits vorhanden sein. 

- `kubectl`
- `docker`
- `helm`

## Get started

Zur Vereinfachung des Prozesses wurde ein _script_ `init.sh` geschrieben.

Beim Ausführen werden die __*Anforderungen*__ installiert. Zunächst werden die vorhandenen _helm charts_ deinstalliert
und anschließend `element` und `synapse` deployed. Alternative zu `synapse` kann diese durch `dendrite` ersetzt werden.

__Es ist zu beachten, dass der default *homeserver* `synapse` ist und die `element` Einstellungen
darauf angepasst worden sind!__ Wenn der `dendrite` als *homeserver* gewünscht ist, sollten der _Name_ und _Port_
dementsprechend in `values.yaml` unter element Ordner angepasst werden.

Zum Start mit dem *default* Einstellungen führe den folgenden Befehl in der Console:

```shell
bash init.sh
```

Der FE und damit `element` kann durch die Adresse [http://localhost:80](http://localhost:80) erreicht werden.

## Anpassung

### FE

Der FE wurde unter `m.homeserver` die `base_url` zusammen mit dem Port und der `server_name` angepasst. Das heißt
wiederum, dass wenn der *homeserver* oder dessen *Port* sich ändern sollten, diese dementsprechend beim FE
angepasst werden müssen.

### BEs

Beide Implementierungen geben den Port `30009` frei.

__(Nur Synapse)__: Die Konfiguration für OpenID-Connect wurde bei `values.yaml` unter `oidc_providers` ergänzt. Diese
wurden aus der [Dokumentation](https://matrix-org.github.io/synapse/latest/openid.html) von `synapse` entnommen.
Hier wird u. a. festgelegt wie der *externe* User im Matrix angelegt wird. Hierzu siehe `values.yaml` unter synapse
Ordner.

### Middleware service

Dieser Service stellt der Admin-Bot dar. Ferner sollte es als PoC für das automatisierte und granulare API-Ansprechen
von dem Matrix-BE dienen.

Der Service ist in Go geschrieben und ist einfach und modular aufgebaut. Dennoch *muss* der Service manuell gestartet
werden. __Der richtige Zeitpunkt ist nach dem erfolgreichen Einloggen der User__. Weiterhin ist es zu beachten, dass
der angesprochene *Dummy User* __@alan:localhost__ hardgecoded ist. Diese *sollte* im Falle einer OIDC-Anbindung geändert
werden.

Bei seiner *initialen* Funktion wird der Service den BE nach einem bestimmten User `@admin:localhost` abfragen und wenn
dieser nicht vorhanden ist, wird diesen angelegt. Wenn der User aber bereits vorhanden sein sollte, wird der Service
sich als solchen beim Matrix-BE einloggen.

Anschließend wird durch den Admin-Bot ein Raum erstellt, den Dummy-User zu diesem Raum hinzugefügt und gleichzeitig die
Raum-Rechte des Users festgelegt. Es ict auch möglich die Rechte nach der Erstellung und User-Zuweisung anpassen. Es muss
aber einen anderen API angesprochen werden (siehe entsprechendes Kommentar in `main.go`).

Der Raum ist ein unidirektionaler Raum, indem nur bestimmte User Nachrichten senden und alle anderen nur diese empfangen
können. Solcher Raum wird als *Infokanal* bezeichnet.
