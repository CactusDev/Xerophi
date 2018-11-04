
import { Model, Column } from "pims";

@Model({
	table: "trusts"
})
export class Trusts {

	@Column({ primary: true })
	public id: string;

	@Column() public channel: string;
	@Column() public platform: string;
	@Column() public userId: string;
}
