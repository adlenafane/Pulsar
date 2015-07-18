# Pulsar

![alt tag](http://corsair.space/pulsar_hud.png)
![alt tag](http://corsair.space/pulsar_ds1.png)

The Pulsar project aim to visually render the mood, trend activity of an Application over time. Demo at http://corsair.space (Mood of itself :))

Installation:
-------------

- InfluxDB
- Frontend (Serve HTTP)

Configuration File:
-------------------

[] Todo

Usage Backend:
--------------

- Emit pulse:

```
curl '127.0.0.1:4000/pulse?token=<APPLICATION_TOKEN>&client_id=<CLIENT_ID>&group_id=<GROUP_ID>'
```
APPLICATION_TOKEN: Autogenerated token for your application so your activity is visible in the right place
CLIENT_ID: Unique identifier for your client (nb: you can send a hash of your id as long as it is unique per client)
GROUP_ID: *Optionnal, group your client belongs to. This is only to visually draw Nebula of users

- Get Statistics

```
curl '127.0.0.1:4000/stats/loads?token=<APPLICATION_TOKEN>'

curl '127.0.0.1:4000/stats?token=<APPLICATION_TOKEN>&groupBy=<AGGREGATE_TIME>&past=<START_TIME>'
```
AGGREGATE_TIME: On which range aggregation shall be displayed (ex: 20min / 1h)
START_TIME: range will be taken on now() - START_TIME (ex: 1d, 12h)

Nb: Be carefull when choosing your representation you may experience trouble if querying too many data (ex: past=10d & groupBy=1s --> 10*24*3600 point in your graph)

- Receive pulse notification (websocket)

```
var sock = new WebSocket("ws://127.0.0.1:4000/sock");
sock.send('{"action":"joingalaxy","data":"33999f3f30f9"}');
```
On first connection you must send the Galaxy you want to join using your APPLICATION_TOKEN.

Usage Frontend
--------------

Host FrontEnd

Philosophy:
-----------

Simple is best

```
GalaxyCluster {
	[]Galaxy {
		[]Nebula {
			[]Pulsar
		}
	}
}
```
Pulsar: Single Client
Nebula: Client Organization/Group
Galaxy: Application


Ex:

Application	: ProductStream
Nebula		: Alkemics
Pulsars		: ANY_STRING_CLIENT_UNIQ;

Author:
-------

```
Samuel RAMOND
@erazor42
```








