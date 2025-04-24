# Cloudflare Go, Update Load Balancer Pools by Endpoint Name
<SCREENSHOT HERE>

This script offers the ability to list all endpoints or specific endpoints by their names within a given Cloudflare Account. It also offers the ability to update the Endpoint's status as to whether or not it is enabled.

If you wish to only list out endpoints with the script, do not send any ENDPOINT environment variables. If you wish to look at a specific ENDPOINT, you can set `CF_ENDPOINT_NAME` and then subsequently `CF_ENDPOINT_ACTION` to `get`, `enable`, or `disable`. The two latter options will update the endpoint in the Cloudflare dashboard and make API calls at the end of scanning all pools.

To use this script, ensure the following environment variables are set in the environment in which the Go program is executed.

| Name       | Description                                                                                     |
|------------|-------------------------------------------------------------------------------------------------|
| CF_ACCT_ID  | This is your Cloudflare Account ID number which can be pulled on the home page of your CF dashboard. |
| CF_API_EMAIL   | This is the email of a user associated with the account that has sufficient permissions.            |
| CF_API_KEY     | This can either be a global key, or you can generate your own token. If you wish to use double check that variable names don't also change to tokens |
| CF_ENDPOINT_NAME     | This is the name of the given endpoint (formerly known as origin) that you wish to update. Only set this if you wish to fetch a specific endpoint or make changes to endpoints, otherwise the script will list all endpoints and their current configuration. |
| CF_ENDPOINT_ACTION     | This the action that will be taken on the endpoint selected by name. Actions supported are `get` which will only show those endpoint's details. `enable` to enable the endpoint. `disable` to disable the endpoint. |