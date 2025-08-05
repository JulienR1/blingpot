import { createFormHook } from "@tanstack/react-form";
import { fieldContext, formContext } from "./form-context";

import { Input } from "./input";
import { Label } from "./label";
import { Combobox } from "./combobox";
import { DatePicker } from "./date-picker";
import { Button } from "./button";
import { FieldError } from "./field-error";

const { useAppForm } = createFormHook({
  fieldComponents: { Input, Combobox, DatePicker },
  formComponents: { Label, Button, FieldError },
  fieldContext,
  formContext,
});

export { useAppForm };
