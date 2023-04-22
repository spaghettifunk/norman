package commander

/*

The commander service is responsible for the following:

* Maintaining global metadata (e.g. configs and schemas) of the system with
the help of Raft which is used as the persistent metadata store.

* Managing other Norman components (brokers and storages)

* Maintaining the mapping of which servers are responsible for which segments.
This mapping is used by the servers to download the portion of the segments that they are responsible for.
This mapping is also used by the broker to decide which servers to route the queries to.

* Serving admin endpoints for viewing, creating, updating, and deleting configs,
which are used to manage and operate the cluster.

* Serving endpoints for segment uploads, which are used in offline data pushes.
They are responsible for initializing real-time consumption and coordination of
persisting real-time segments into the segment store periodically.

* Undertaking other management activities such as managing retention of segments, validations.

*/
