module Todo exposing (..)
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)

-- MAIN
main : Program() Model Msg
main =
  Browser.sandbox { init = init, update = update, view = view }

-- MODEL
type alias Model = { title : String, description : String, completed : Bool}

init : Model
init = Model "title" "description" False

-- UPDATE
type Msg =  Complete | Incomplete

update : Msg -> Model -> Model
update msg model =
  case msg of
    Complete -> {model | completed = True}
    Incomplete -> {model | completed = False}

-- VIEW
view : Model -> Html Msg
view model =
  div [] [
    div [class "header"] [ h1 [] [text "Todo List"]],
    div [class "task"] [
        h2 [] [text "make dinner"], text "make spaghetti and prepare wine", viewIcon model
        ] 
  ]

viewIcon : Model -> Html Msg
viewIcon model =
    let 
        iconType = if model.completed then "select_check_box" else "check_box_outline_blank"
        msg = if model.completed then Incomplete else Complete
    in

    div [class "checkbox-icon"] [ span [class "material-icons", onClick msg] [ text iconType] ]