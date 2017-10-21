
import * as Hapi from "hapi";
import * as Boom from "boom";

export class ConfigController {

	public async getConfig(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const name = request.params["name"];

		// TODO: Make this actually pull from a database and display information
		const response: Config = {
			repeat: {
					disabled: false,
					onlyLive: true,
					defaultMinimum: 60
				},
				events: {
					follow: {
						message: "Thanks for following, %USER%!",
						enabled: true
					},
					subscribe: {
						message: "Thanks for subscribing, %USER%!",
						enabled: true
					},
					host: {
						message: "Thanks for hosting the channel, %USER% (%VIEWERS% viewers)",
						enabled: false
					},
					join: {
						message: "Welcome, %USER%",
						enabled: false
					},
					leave: {
						message: "Bye, %USER%",
						enabled: false
					}
				},
				whitelistedURLs: [],
				spam: {
					allowUrls: {
						action: "purge",
						value: false,
						warnings: 1
					},
					maxCaps: {
						action: "purge",
						value: 10,
						warnings: 1
					},
					maxEmoji: {
						action: "purge",
						value: 2,
						warnings: 1
					},
					keywords: {
						blacklist: [],
						whitelist: []
					}
				}			
		};
		reply(response);
	}
}