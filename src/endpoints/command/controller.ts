
import * as Hapi from "hapi";
import * as Boom from "boom";

export class CommandController {

	public async getCommand(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const name = request.params["name"];
		const channel = request.params["channel"];

		// TODO: Make this actually pull from a database and display information
		const response: Command = {
			name: name,
			channel: channel,
			response: [
				{
					type: "message",
					action: false,
					text: [
						{
							type: "text",
							data: "Hello!"
						}
					]
				}
			],
			count: 0,
			enabled: true,
			restrictions: {
				service: [],
				role: "user"
			}
		};
		reply(response);
	}
}