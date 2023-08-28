# New Infinite Scale Releas Info

If there is a new Infinite Scale release on the way, there are some things to know respectively to check:

**Relevant for docs:**

* The docs relevant releasing process in the `ocis` repo works the following:
  * Based on a release tag like `v4.0.0`, a new `stable-x` branch is created.
  * A new and empty `docs-stable-x` branch is created.
  * The pipeline of the `stable-x` branch gets adapted to write to the `docs-stable-x` branch.
  * The `stable-x` branch gets added to the nightly cron jobs in drone which populates new and changed content to `docs-stable-x`.
  * Any merge in the `stable-x` branch also triggers the pipeline as usual.

* **IMPORTANT** For any changes necessary that effect the output of the `docs-stable-x`, these must be done in `stable-x` and take effect in `docs-stable-x` when the `stable-x` pipeline ran. Manual commits to `docs-stable-x` will get overwritten at least daily.

* **IMPORTANT**: Any changes neccessary in `stable-x` are based on a 6-eye principle and need 2 approvals therefore.

**Actions for docs:**

* Separate the `antora.yml` updates into the ocis and Helm Chart part. This eases testing.

* You can prepare the changes, but testing can only start when the `release tag` and the `docs-stable-x` branch is available.
* For ocis, the attributes used for assembling paths and printed names in `antora.yml` need to be adapted accordingly:  
Path building: `service_tab_x`, `compose_tab_x`  
Naming: `service_tab_x_tab_text` and `compose_tab_x_tab_text`

* For Helm Charts, the attributes used for assembling paths and printed names in `antora.yml` need to be adapted accordingly:  
Path building: `helm_tab_x`  
Naming: `helm_tab_x_tab_text`

* Make a build and check validity of the
  * tabs + content in the envvar tab in services and
  * links (paths) used in eg. `deployment/container/orchestration/orchestration.html#docker-compose-examples`.
