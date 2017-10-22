
import * as Hapi from "hapi";
import * as Boom from "boom";

import { MongoHandler } from "../../mongo";

import { defaultScopes } from "../../authorization/scopes";
import { hash } from "../../util";

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

	public async createUser(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params.channel;

		if (!request.payload || !request.payload.password) {
			return reply(Boom.badData("Must supply password"));
		}

		const scopes = !!request.payload.scopes ? request.payload.scopes : defaultScopes
		const hashed = await hash(request.payload.password);

		const user = await this.mongo.createUser(channel, hashed, scopes);
		if (!user) {
			return reply(Boom.conflict("User already exists"));
		}
		return reply({
			created: true,
			scopes,
			username: user.username
		}).code(201);
	}

	public async removeUser(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params.channel;

		const deleted = await this.mongo.softDeleteUser(channel);
		if (!deleted) {
			return reply(Boom.notFound("Invalid user"));
		}
		return reply({}).code(204);
	}
}