
import { Config } from "../../config";
import { AbstractEndpoint } from "..";
import { Web } from "../../web";
import { CommandController } from "./controller";

export class CommandRoute extends AbstractEndpoint {

	private controller = new CommandController();

	public async init() {
		this.web.instance.route({
			method: "GET",
			path: "/{channel}/command/{command}",
			config: {
				handler: (request, reply) => this.controller.getCommand(request, reply)
			}
		});
	}
}