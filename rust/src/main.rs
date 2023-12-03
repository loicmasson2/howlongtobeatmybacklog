use serde::Deserialize;
use yew::prelude::*;
use gloo_net::{http::Request, Error};
use gloo_console::log;
use wasm_bindgen::JsValue;

// url that can be called
// http://127.0.0.1:8080/get-games-from-db/76561198000800114
#[derive(Clone, PartialEq, Deserialize)]
struct Data {
    steam_id: String,
    games: Vec<Game>,
}

#[derive(Clone, PartialEq, Deserialize)]
struct Response {
    data: Data,
}

#[derive(Clone, PartialEq, Deserialize)]
struct Game {
    appid: usize,
    img_icon_url: String,
    name: String,
    playtime_disconnected: usize,
    playtime_forever: usize,
    playtime_linux_forever: usize,
    playtime_mac_forever: usize,
    playtime_windows_forever: usize,
    rtime_last_played: usize,
}

#[derive(Properties, PartialEq)]
struct VideoGamesListProps {
    videogames: Vec<Game>,
    on_click: Callback<Game>,
}

#[function_component(VideoGamesList)]
fn videogames_list(crate::VideoGamesListProps { videogames, on_click }: &crate::VideoGamesListProps) -> Html {
    let on_click = on_click.clone();
    videogames.iter().map(|videogame| {
        let on_videogame_select = {
            let on_click = on_click.clone();
            let videogame = videogame.clone();
            Callback::from(move |_| {
                on_click.emit(videogame.clone())
            })
        };
        html! {
            <div style="display: flex; flex-direction: row; align-items: center; gap: 8px;">
                <img width="24px" height="24px" src={format!("http://media.steampowered.com/steamcommunity/public/images/apps/{}/{}.jpg", videogame.appid, videogame.img_icon_url)} alt="videogame thumbnail" />
                <p style="margin-top: 8px; margin-bottom: 8px;" key={videogame.appid} onclick={on_videogame_select}>{format!("{}: {}", videogame.name, videogame.playtime_forever)}</p>
            </div>
    }
    }).collect()
}

#[derive(Properties, PartialEq)]
struct VideoGameDetailsProps {
    videogame: Game,
}

#[function_component(VideoGameDetails)]
fn videogame_details(crate::VideoGameDetailsProps { videogame }: &crate::VideoGameDetailsProps) -> Html {
    html! {
        <div>
            <h3>{ videogame.name.clone() }</h3>
            <img src={format!("https://cdn.akamai.steamstatic.com/steam/apps/{}/header.jpg?t=1701200506", videogame.appid)} alt="video thumbnail" />
        </div>
    }
}


#[function_component(App)]
fn app() -> Html {
    let videogames = use_state(|| vec![]);
    let error: UseStateHandle<Option<Error>> = use_state(|| None);
    {
        let videogames = videogames.clone();
        use_effect_with((), move |_| {
            let videogames = videogames.clone();
            wasm_bindgen_futures::spawn_local(async move {
                let fetched_videogames: Response = Request::get("http://127.0.0.1:8080/get-games-from-db/76561198000800114")
                    .send()
                    .await
                    .unwrap()
                    .json()
                    .await
                    .unwrap();
                videogames.set(fetched_videogames.data.games);
            });
            || ()
        });
    }

    let selected_videogame = use_state(|| None);
    let on_videogame_select = {
        let selected_videogame = selected_videogame.clone();
        Callback::from(move |videogame: Game| {
            selected_videogame.set(Some(videogame))
        })
    };

    let details = selected_videogame.as_ref().map(|videogame| html! {
        <VideoGameDetails videogame={videogame.clone()} />
    });

    html! {
        <div>
           <h1>{ "How Long To Beat My Backlog" }</h1>
           <div style="display: flex; flex-direction: row;">
              <div style="width: 50%;">
                 <h3>{"Games to watch"}</h3>
                 <VideoGamesList videogames={(*videogames).clone()} on_click={on_videogame_select.clone()} />
              </div>
              <div style="width: 50%;">
                 { for details }
              </div>
           </div>
        </div>
    }
}

fn main() {
    // wasm_logger::init(wasm_logger::Config::default());
    yew::Renderer::<App>::new().render();
}