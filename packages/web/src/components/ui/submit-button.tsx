import type { ComponentProps } from "react";
import { Button } from "./button";
import { Loader2Icon } from "lucide-react";

interface SubmitButtonProps extends ComponentProps<typeof Button> {
  canSubmit: boolean;
  isSubmitting: boolean;
}

export function SubmitButton({
  canSubmit,
  isSubmitting,
  ...props
}: SubmitButtonProps) {
  return (
    <Button {...props} type="submit" disabled={props.disabled || !canSubmit}>
      {isSubmitting && <Loader2Icon className="animate-spin" />}
      {props.children}
    </Button>
  );
}
