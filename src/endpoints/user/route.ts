
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { UserController } from "./controller";

import { ADD_USER, REMOVE_USER } from "../../authorization/scopes";

export class UserRoute extends AbstractEndpoint {

	private controller: UserController;

	public async init() {
		this.controller = new UserController(this.web.mongo);

		this.web.instance.route([
			{
				method: "GET",
				path: "/user/login",
				config: {
					handler: (request, reply) => this.controller.attemptLogin(request, reply, this.config.authentication.secret),
					auth: false
				}
			},
			{
				method: "POST",
				path: "/user/create",
				config: {
					handler: (request, reply) => this.controller.createUser(request, reply),
					auth: {
						scope: ADD_USER
					}
				}
			},
			{
				method: "DELETE",
				path: "/user/{channel}",
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