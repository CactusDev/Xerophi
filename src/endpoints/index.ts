
import { Endpoint, Get, Describe } from "../annotation/endpoint";
import { ChannelCache, ChannelBucket, RequestVerb } from "../internal";

@Endpoint({
	cache: ChannelCache,
	bucket: ChannelBucket
})
export class ChannelEndpoint {

	@Describe("/{name}")
	@Get({
		authorization: "channel:view"
	})
	public async getChannel(): Promise<object> {

	}
}
