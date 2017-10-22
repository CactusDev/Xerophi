
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { CommandController } from "./controller";

export class CommandRoute extends AbstractEndpoint {

	private controller: CommandController;

	public async init() {
		this.controller = new CommandController(this.web.mongo);

		this.web.instance.route({
			method: "GET",
			path: "/{channel}/command/{command}",
			config: {
				handler: (request, reply) => this.controller.getCommand(request, reply),
				auth: false
			}
		});

		this.web.instance.route({
			method: "POST",
			path: "/{channel}/command/{command}",
			config: {
				handler: (request, reply) => this.controller.createCommand(request, reply)
			}
		});

		this.web.instance.route({
			method: "PATCH",
			path: "/{channel}/command/{command}",
			config: {
				handler: (request, reply) => this.controller.updateCommand(request, reply)
			}
		});
	}
}