
import { Endpoint, Get, Delete, Post, Patch, Describe } from "../annotation/endpoint";
import { ChannelCache, ChannelBucket } from "../internal";

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

	@Describe("/new")
	@Post()
	public async createChannel(): Promise<object> {

	}

	@Describe("/{name}")
	@Delete({
		authorization: "channel:delete"
	})
	public async deleteChannel(): Promise<object> {

	}

	@Describe("/{name}")
	@Patch({
		authorization: "channel:update"
	})
	public async updateChannel(): Promise<object> {

	}
}
