
import { Model, Column } from "pims";

@Model({
	table: "permits"
})
export class Permits {

	@Column({ primary: true })
	public id: string;

	@Column() public channel: string;
	@Column() public platform: string;
	@Column() public userId: string;
	@Column() public offenceType: string;
	@Column() public time: string;
	@Column() public reason: string;
}
