import { Injectable } from '@angular/core';
import { FieldBase } from './field-base';
import { FieldTextbox } from './field-textbox';

import { of } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class FieldService {

  getFields() {
    const test: FieldBase<string>[] = [
      new FieldTextbox({
        key: 'testKey',
        label: 'testLabel',
        value: 'testValue',
        required: true,
      })
    ]; 

    return of(test);
  };
}
