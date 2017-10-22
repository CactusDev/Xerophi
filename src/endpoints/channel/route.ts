
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { ChannelController } from "./controller";

import { ADD_USER, REMOVE_USER } from "../../authorization/scopes";

export class ChannelRoute extends AbstractEndpoint {

	private controller: ChannelController;

	public async init() {
		this.controller = new ChannelController(this.web.mongo);

		this.web.instance.route([
			{
				method: "GET",
				path: "/channel/{channel}",
				config: {
					handler: (request, reply) => this.controller.getChannel(request, reply),
					auth: false
				}
			},
			{
				method: "GET",
				path: "/channel/{channel}/{service}",
				config: {
					handler: (request, reply) => this.controller.getService(request, reply),
					auth: false
				}
			},
			{
				method: "POST",
				path: "/channel/{channel}",
				config: {
					handler: (request, reply) => this.controller.createUser(request, reply),
					auth: {
						scope: ADD_USER
					}
				}
			},
			{
				method: "DELETE",
				path: "/channel/{channel}",
				config: {
					handler: (request, reply) => this.controller.removeUser(request, reply),
					auth: {
						scope: REMOVE_USER
					}
				}
			}
		]);
	}
}