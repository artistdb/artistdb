import { Injectable } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';

import { FieldBase } from './field-base';

@Injectable({
  providedIn: 'root'
})
export class FieldControlService {
  constructor() { }

  toFormGroup(fields: FieldBase<string>[]) {
    const group: any = {};

    fields.forEach(field => {
      group[field.key] = field.required ? new FormControl(field.value || '', Validators.required)
        : new FormControl(field.value || '');
    });
    return new FormGroup(group);
  }
}
