import { FieldBase } from "./field-base";

export class FieldTextbox extends FieldBase<string> {
    override controlType = 'textbox';
}