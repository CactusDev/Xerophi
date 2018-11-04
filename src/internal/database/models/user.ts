
import { Model, Column } from "pims";

@Model({
	table: "users"
})
export class Users {
	@Column({ primary: true })
	public id: string;

	@Column() public channels: string[];
	@Column() public token: string;
	@Column() public deletedAt: number;
	@Column() public passwordHash: string;
	@Column() public adminFlags: number;
	@Column() public timezone: string;
	@Column() public pointsName: string;
	@Column() public customVariables: {[name: string]: string};
}
