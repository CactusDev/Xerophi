
import * as Hapi from "hapi";
import * as Boom from "boom";

const aguid = require("aguid");

import { Authorization } from "../../authorization";
import { defaultScopes } from "../../authorization/scopes";

import { compare, hash } from "../../util";
import { MongoHandler } from "../../mongo";

export class UserController {

	constructor(private mongo: MongoHandler) {
	}

	public async attemptLogin(request: Hapi.Request, reply: Hapi.ReplyNoContinue, key: string) {
		const user = request.headers["user"];
		const password = request.headers["password"];
		
		const dbUser = await this.mongo.getUser(user);

		if (!dbUser) {
			return reply(Boom.notFound("Invalid user"));
		}

		if (await compare(password, dbUser.passwordHash)) {
			// Valid user, let them exist!
			const session: any = {
				valid: true,
				id: aguid(),
				scopes: dbUser.scopes,
				exp: new Date().getTime() + 60 * 60 * 1000 // This will expire in one hour
			};
			return reply({
				accepted: true,
				jwt: await Authorization.give(session, key)
			});
		}
		// Invalid login
		reply(Boom.unauthorized());
	}

	public async createUser(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		if (!request.payload || !request.payload.password || !request.payload.username) {
			return reply(Boom.badData("Must supply password"));
		}

		const scopes = !!request.payload.scopes ? request.payload.scopes : defaultScopes
		const channel = request.payload.username;
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