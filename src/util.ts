
const validRoles = ["banned", "user", "subscriber", "moderator", "owner"];

export function isValidRole(role: string): boolean {
	return validRoles.indexOf(role) > -1;
}
