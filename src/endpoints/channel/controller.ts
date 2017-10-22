
import * as Hapi from "hapi";
import * as Boom from "boom";

export class ChannelController {

	public async getChannel(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params["channel"];
		const service = request.params["service"];

		// TODO: Make this actually pull from a database and display information
		reply(!!service ? this.hasService(channel, service) : this.noService(channel)).header("Authorization", request.headers.authorization);
	}

	private hasService(channel: string, service: string): Channel {
		return {
			repeats: [],
			username: "0x01",
			service: "twitch",
			uuid: "",
			trusts: ["Innectic", "2Cubed", "ParadigmShift3d"],
			permits: [],
			chatters: {
				Innectic: {
					points: 100
				}
			},
			config: {
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
			}
		};
	}

	private noService(channel: string): User {
		return {
			username: channel,
			deletedAt: null,
			uuid: "",
			passwordHash: "",
			channels: [],
			scopes: [],
			commands: []
		}
	}
}