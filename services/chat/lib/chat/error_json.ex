defmodule Chat.ErrorJSON do
  def render("404.json", _assigns), do: %{errors: %{detail: "Not Found"}}
  def render("500.json", _assigns), do: %{errors: %{detail: "Internal Server Error"}}

  def render(template, _assigns),
    do: %{errors: %{detail: Phoenix.Controller.status_message_from_template(template)}}
end
