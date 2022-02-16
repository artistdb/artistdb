import { Component, Input, OnInit } from '@angular/core';
import { FormGroup } from '@angular/forms';

import { FieldBase } from './field-base';
import { FieldControlService } from './field-control.service';

@Component({
  selector: 'app-dynamic-form',
  templateUrl: './dynamic-form.component.html',
  providers: [ FieldControlService ]
})
export class DynamicFormComponent implements OnInit {
  @Input() fields: FieldBase<string>[] | null = [];
  form!: FormGroup;
  payLoad = '';

  constructor(private fcs: FieldControlService) {}

  ngOnInit() {
    this.form = this.fcs.toFormGroup(this.fields as FieldBase<string>[]);
  }

  onSubmit() {
    this.payLoad = JSON.stringify(this.form.getRawValue());
  }
}