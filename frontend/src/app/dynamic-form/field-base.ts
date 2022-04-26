export class FieldBase<T> {
    value: T|undefined;
    key: string;
    label: string;
    required: boolean;
    controlType: string;
    type: string;

  constructor(options: {
    value?: T;
    key?: string;
    label?: string;
    required?: boolean;
    controlType?: string;
    type?: string;
  } = {}) {
  this.value = options.value;
  this.key = options.key || '';
  this.label = options.label || '';
  this.required = !!options.required;
  this.controlType = options.controlType || '';
  this.type = options.type || '';
}
}