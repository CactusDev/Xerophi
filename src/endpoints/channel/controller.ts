
import * as Hapi from "hapi";
import * as Boom from "boom";

import { MongoHandler } from "../../mongo";

import { defaultScopes } from "../../authorization/scopes";

interface APIUser {
	username: string;
	password: string;
	scopes: string;
}

const usersKeys = ["username", "password", "scopes"];

export class ChannelController {

	constructor(public mongo: MongoHandler) {
	}

	private async isValidUser(user: any): Promise<boolean> {
		for (let key of usersKeys) {
			if (!user[key]) return false;
		}
		return true;
	}

	public async getService(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params.channel;
		const service = request.params.service;

		const dbChannel = await this.mongo.getUser(channel);
		if (!dbChannel || dbChannel.deletedAt) {
			return reply(Boom.notFound("Invalid user"));
		}

		if (dbChannel.channels.length === 0) {
			return reply(Boom.notFound("Channel doesn't have service"));
		}
		// See if we have the service
		dbChannel.channels.forEach(async serviceChannel => {
			if (serviceChannel.service === service) {
				// This is the channel we're looking for, remove all the special data
				delete serviceChannel.auth;
				reply(serviceChannel);
			}
		});
	}

	public async getChannel(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params.channel;

		const dbChannel = await this.mongo.getUser(channel);
		if (!dbChannel || dbChannel.deletedAt) {
			return reply(Boom.notFound("Invalid user"));
		}

		delete dbChannel.passwordHash;
		delete dbChannel.deletedAt;

		for (let service of dbChannel.channels) {
			delete service.auth;
		}

		reply(dbChannel);
	}
}