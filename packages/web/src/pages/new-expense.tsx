import { createRoute, Link } from "@tanstack/react-router";
import { root } from "./root";
import z from "zod";
import { useSuspenseQuery } from "@tanstack/react-query";
import { profileQuery, profilesQuery } from "@/stores/profile";
import { useAppForm } from "@/components/ui/form";
import { useCreate } from "@/stores/expense";
import { categoriesQuery } from "@/stores/category";
import { useMemo } from "react";

function NewExpense() {
  const create = useCreate();
  const profile = useSuspenseQuery(profileQuery);
  const profiles = useSuspenseQuery(profilesQuery);
  const categories = useSuspenseQuery(categoriesQuery);

  const defaultValues = useMemo(
    () => ({
      label: "",
      amount: "",
      timestamp: new Date(),
      spenderId: profile.data?.sub ?? "",
      categoryId: null,
    }),
    [profile.data?.sub],
  );

  const form = useAppForm({
    defaultValues: defaultValues as Omit<typeof defaultValues, "categoryId"> & {
      categoryId: number | null;
    },
    validators: {
      onChange: z.object({
        label: z
          .string({ error: "Saisir une description" })
          .min(1, { error: "Saisir une description" }),
        amount: z.coerce
          .number<string>({ error: "Saisir un montant" })
          .positive({ error: "Le montant doit être positif" }),
        timestamp: z.date({ error: "Saisir la date de la transaction" }),
        spenderId: z
          .string({
            error: "Saisir la personne ayant effectué la transaction",
          })
          .refine((sub) => profiles.data.find((p) => p.sub === sub), {
            error: "Cette personne n'est pas disponible",
          }),
        categoryId: z
          .number({ error: "Saisir la catégorie de la transaction" })
          .refine((id) => categories.data.find((c) => c.id === id), {
            error: "Cette catégorie n'est pas disponible",
          })
          .nullable(),
      }),
    },
    onSubmit: ({ value }) => create(value),
  });

  return (
    <>
      <Link to="/">go to index</Link>

      <div>
        <form
          className="max-w-sm flex flex-col justify-center gap-2 py-2"
          onSubmit={(e) => {
            e.preventDefault();
            form.handleSubmit();
          }}
        >
          <form.AppField name="label">
            {(field) => (
              <>
                <form.Label>Description</form.Label>
                <field.Input />
                <form.FieldError />
              </>
            )}
          </form.AppField>

          <form.AppField name="amount">
            {(field) => (
              <>
                <form.Label>Montant</form.Label>
                <field.Input type="number" step="0.01" />
                <form.FieldError />
              </>
            )}
          </form.AppField>

          <form.AppField name="categoryId">
            {(field) => (
              <>
                <form.Label>Catégorie</form.Label>
                <field.Combobox
                  placeholder="Sélectionner une catégorie"
                  search={{
                    placeholder: "Trouver une catégorie",
                    empty: "Aucune catégorie ne correspond",
                  }}
                  options={categories.data.map((category) => ({
                    label: category.label,
                    value: category.id,
                  }))}
                />
                <form.FieldError />
              </>
            )}
          </form.AppField>

          <form.AppField name="timestamp">
            {(field) => (
              <>
                <form.Label>Date</form.Label>
                <field.DatePicker placeholder="Sélectionner une date" />
                <form.FieldError />
              </>
            )}
          </form.AppField>

          <form.AppField name="spenderId">
            {(field) => (
              <>
                <form.Label>Personne</form.Label>
                <field.Combobox
                  placeholder="Sélectionner une personne"
                  search={{
                    placeholder: "Trouver une personne",
                    empty: "Aucune personne ne correspond",
                  }}
                  options={profiles.data.map((profile) => ({
                    label: profile.firstName + " " + profile.lastName,
                    value: profile.sub,
                  }))}
                />
                <form.FieldError />
              </>
            )}
          </form.AppField>

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
            children={([canSubmit, isSubmitting]) => (
              <form.SubmitButton
                canSubmit={canSubmit}
                isSubmitting={isSubmitting}
              >
                {isSubmitting ? "Enregistrement en cours" : "Enregistrer"}
              </form.SubmitButton>
            )}
          />
        </form>
      </div>
    </>
  );
}

export const newExpense = createRoute({
  getParentRoute: () => root,
  path: "/new",
  component: NewExpense,
});
