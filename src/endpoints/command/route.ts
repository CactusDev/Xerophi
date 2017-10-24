
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
					handler: this.controller.getCommand,
					auth: false
				}
			},
			{
				method: "POST",
				path: "/{channel}/command/{command}",
				config: {
					handler: this.controller.createCommand
				}
			},
			{
				method: "PATCH",
				path: "/{channel}/command/{command}",
				config: {
					handler: this.controller.updateCommand
				}
			},
			{
				method: "DELETE",
				path: "/{channel}/command/{command}",
				config: {
					handler: this.controller.softDeleteCommand
				}
			}
		]);
	}
}