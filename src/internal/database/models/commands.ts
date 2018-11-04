
import { Column, Model } from "pims";

@Model({
	table: "commands"
})
export class Commands {

	@Column({ primary: true })
	public id: string;

	@Column() public channel: string;
	@Column() public name: string;
	@Column() public count: number;
	@Column() public enabled: boolean;
	@Column() public response: CommandResponse;
	@Column() public token: string;
}
