import { useFieldContext } from "./form-context";

export function FieldError() {
  const field = useFieldContext();

  return (
    field.state.meta.isValid === false && (
      <em role="alert">
        {field.state.meta.errors
          .map((e) => (typeof e === "object" ? e.message : (e ?? "")))
          .join(", ")}
      </em>
    )
  );
}
