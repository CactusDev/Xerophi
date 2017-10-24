
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { CommandController } from "./controller";

export class CommandRoute extends AbstractEndpoint {

	private controller: CommandController;

	public async init() {
		this.controller = new CommandController(this.web.mongo);

		this.web.instance.route([
			{
				method: "GET",
				path: "/{channel}/command/{command}",
				config: {
					handler: (request, reply) => this.controller.getCommand(request, reply),
					auth: false
				}
			},
			{
				method: "POST",
				path: "/{channel}/command/{command}",
				config: {
					handler: (request, reply) => this.controller.createCommand(request, reply)
				}
			},
			{
				method: "PATCH",
				path: "/{channel}/command/{command}",
				config: {
					handler: (request, reply) => this.controller.updateCommand(request, reply)
				}
			},
			{
				method: "DELETE",
				path: "/{channel}/command/{command}",
				config: {
					handler: (request, reply) => this.controller.softDeleteCommand(request, reply)
				}
			}
		]);
	}
}