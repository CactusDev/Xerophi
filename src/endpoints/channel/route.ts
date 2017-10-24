
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
					handler: this.controller.getChannel,
					auth: false
				}
			},
			{
				method: "GET",
				path: "/channel/{channel}/{service}",
				config: {
					handler: this.controller.getService,
					auth: false
				}
			},
			{
				method: "POST",
				path: "/channel/{channel}",
				config: {
					handler: this.controller.createUser,
					auth: {
						scope: ADD_USER
					}
				}
			},
			{
				method: "DELETE",
				path: "/channel/{channel}",
				config: {
					handler: this.controller.removeUser,
					auth: {
						scope: REMOVE_USER
					}
				}
			}
		]);
	}
}