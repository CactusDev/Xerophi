
import * as Hapi from "hapi";
import * as Boom from "boom";

const argon2 = require("argon2");
const aguid = require("aguid");

import { Authorization } from "../../authorization";

const logins: any = {
	"0x01": "$argon2i$v=19$m=4096,t=3,p=1$2WC62WsiICG2rnfToHtPpw$Kpf6d2N+qLmhCJgZKYSUn1hDMIwUUbejzGpkcPGNKwE"
}

// @Temp
const userScopes: {[name: string]: string[]} = {
	"0x01": [
		"user:basic:auth",
		"user:command:create",
		"user:command:delete",
		"user:command:edit",
		"user:quote:create",
		"user:quote:delete"
	],
	test: [
		"admin:full"
	]
}


export class LoginController {

	private async compare(password: string, hash: string): Promise<boolean> {
		try {
			return await argon2.verify(hash, password);
		} catch (e) {
			return false;
		}
	}

	private async hash(password: string): Promise<string> {
		try {
			return await argon2.hash(password);
		} catch (e) {
			return null;
		}
	}

	public async attemptLogin(request: Hapi.Request, reply: Hapi.ReplyNoContinue, key: string) {
		const user = request.headers["user"];
		const password = request.headers["password"];
		if (!logins[user]) {
			return; // TODO: error here
		}

		if (await this.compare(password, logins[user])) {
			// Valid user, let them exist!
			const session: any = {
				valid: true,
				id: aguid(),
				scopes: userScopes[user],
				exp: new Date().getTime() + 60 * 60 * 1000 // This will expire in one hour
			};
			reply({
				accepted: true,
				jwt: await Authorization.give(session, key)
			});
			return;
		}
		// Invalid login
		reply(Boom.unauthorized());
	}
}