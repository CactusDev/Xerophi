
import * as Hapi from "hapi";
import * as Boom from "boom";
import { isValidRole } from "../../util";

import { MongoHandler } from "../../mongo";

export class CommandController {

	constructor(private mongo: MongoHandler) {

	}

	public async validResponse(response: CactusMessagePacket[]): Promise<boolean> {
		return true; // TODO
	}

	public async getCommand(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const name = request.params["command"];
		const channel = request.params["channel"];

		const command = await this.mongo.getCommand(channel, name);
		if (!command) {
			return reply(Boom.notFound("Invalid command"));
		}
		return reply(command);
	}

	public async createCommand(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params["channel"];
		const name = request.params["command"];

		if (!request.payload || !request.payload.response || !request.payload.role) {
			return reply(Boom.badData("Must supply a payload with role & response"));
		}

		const response: CactusMessagePacket[] = !!request.payload && !!request.payload.response ? request.payload.response : null;
		if (!response || !await this.validResponse(response)) {
			return reply(Boom.badData("Invalid response"));
		}

		const role: string = request.payload.role;
		if (!isValidRole(role)) {
			return reply(Boom.badData("Invalid role"));
		}

		// Everything is valid, so lets make a command!
		const created = await this.mongo.createCommand(channel, name, response, role);
		if (!created) {
			return reply(Boom.conflict("Command already exists"));
		}
		// Created, display a nice message
		return reply({
			created: true,
			name
		}).code(201);
	}

	public async updateCommand(request: Hapi.Request, reply: Hapi.ReplyNoContinue) {
		const channel = request.params["channel"];
		const command = request.params["command"];
		
		if (!request.payload || !(request.payload.role ||
			request.payload.response || request.payload.name || request.payload.enabled)) {
			return reply(Boom.badData("Must supply something to update."));
		}

		let updated = false;
		// Since we have all the data, we just need to actually update the attributes.
		if (request.payload.role) {
			// Validate the role
			if (!isValidRole(request.payload.role)) {
				return reply(Boom.badData("Invalid role"));
			}
			updated = await this.mongo.commandEditRestrict(request.payload.role, command, channel);			
		} else if (request.payload.response) {
			updated = await this.mongo.commandEditResponse(request.payload.response, command, channel);
		} else if (request.payload.name) {
			updated = await this.mongo.commandEditName(request.payload.name, command, channel);
		} else if (request.payload.enabled) {
			updated = await this.mongo.commandEditEnabled(request.payload.enabled, command, channel);
		} else {
			return reply(Boom.badData("Invalid attribute"));
		}
		if (!updated) {
			return reply(Boom.notFound("Invalid command."));
		}
		return reply({}).code(204);
	}
}