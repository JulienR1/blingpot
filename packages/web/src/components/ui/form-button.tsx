import type { ComponentProps } from "react";
import { Button } from "./button";
import { Loader2Icon } from "lucide-react";

interface FormButtonProps extends ComponentProps<typeof Button> {
  canSubmit: boolean;
  isSubmitting: boolean;
}

export function FormButton({
  canSubmit,
  isSubmitting,
  ...props
}: FormButtonProps) {
  return (
    <Button {...props} disabled={props.disabled || canSubmit}>
      {isSubmitting && <Loader2Icon className="animate-spin" />}
      {props.children}
    </Button>
  );
}
