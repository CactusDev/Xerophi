
const moment = require("moment-strftime");

const validRoles = ["banned", "user", "subscriber", "moderator", "owner"];

export function isValidRole(role: string): boolean {
	return validRoles.indexOf(role) > -1;
}

export function theTime(): string {
	return moment().strftime("%a %b %d %H:%M:%S %Y");
}
