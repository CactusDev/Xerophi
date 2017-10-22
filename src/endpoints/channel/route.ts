
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { ChannelController } from "./controller";

export class ChannelRoute extends AbstractEndpoint {

	private controller = new ChannelController();

	public async init() {
		this.web.instance.route({
			method: "GET",
			path: "/channel/{channel}",
			config: {
				handler: (request, reply) => this.controller.getChannel(request, reply),
				auth: "jwt"
			}
		});

		this.web.instance.route({
			method: "GET",
			path: "/channel/{channel}/{service}",
			config: {
				handler: (request, reply) => this.controller.getChannel(request, reply),
				auth: "jwt"
			}
		});
	}
}