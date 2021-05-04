kernel.org Transparency Log Monitor
===================================

This project is a monitor for the kernel.org git push-cert transparency log
which is published. It contains a push certificate signed by the commiter and
can allow external parties to verify that commits has been pushed by a certified
develooper.

Currently this project pulls the pgpkeys from kernel.org and attempts to
validate any found signatures towards this keyring.

One can also ask the service if a given revision has been seen by the log by
utilizing the `/api/revision/:commit` endpoint.
