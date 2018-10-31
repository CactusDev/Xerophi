
import { Endpoint } from "../annotation/endpoint";
import { ChannelCache, ChannelBucket, RequestVerb } from "../internal";

@Endpoint({
	cache: ChannelCache,
	bucket: ChannelBucket,
	authentication: {
		[RequestVerb.Post]:   "channel:create",
		[RequestVerb.Delete]: "channel:delete",
		[RequestVerb.Patch]:  "channel:edit"
	}
})
export class ChannelEndpoint {

}
