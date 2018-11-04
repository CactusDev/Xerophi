
import { Model, Column } from "pims";

@Model({
	table: "points"
})
export class Points {

	@Column({ primary: true })
	public id: string;

	@Column() public channel: string;
	@Column() public values: {[user: string]: number};
}
