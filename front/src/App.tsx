import React, { useEffect, useRef, useState } from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import "./App.css";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:4444/";

const Search = ({
  thinking,
  onInputChange,
  onCheckBoxChange,
  onAbortClick,
  onSearchClick,
  onKeyDown,
  inputRef,
}: {
  thinking: boolean;
  onInputChange: (val: string) => void;
  onCheckBoxChange: (val: boolean) => void;
  onAbortClick: () => void;
  onSearchClick: () => void;
  onKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => void;
  inputRef: React.Ref<HTMLInputElement>;
}) => {
  return (
    <div className="search-container">
      <div className="info">{thinking ? "Thinking..." : ""}</div>
      <div className="search">
        <input
          type="checkbox"
          name="Extract?"
          onChange={(e) => onCheckBoxChange(e.target.checked)}
        />
        <input
          ref={inputRef}
          className="input"
          disabled={thinking}
          type="search"
          onKeyDown={onKeyDown}
          onChange={(e) => onInputChange(e.target.value)}
        />
        <button onClick={onSearchClick}>Search</button>
        <button disabled={!thinking} onClick={onAbortClick}>
          Abort
        </button>
      </div>
    </div>
  );
};

const createMarkdownElement = (data: string) => {
  return (
    <Markdown remarkPlugins={[remarkGfm]} className="message">
      {data}
    </Markdown>
  );
};

let controller: AbortController | null = null;
function App() {
  const [query, setQuery] = useState<string>("");
  const [thinking, setThinking] = useState<boolean>(false);
  const [messageElements, setMessageElements] = useState<
    Array<{ answer: React.ReactNode; time?: string; question: React.ReactNode }>
  >([]);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [currentAnswer, setCurrentAnswer] = useState<string>("");
  const [extract, setExtract] = useState<boolean>(false);

  const focusInput = () => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  useEffect(() => {
    focusInput();
  }, []);

  const removeThinkTag = (data: string) =>
    data.replace("<think>", "").replace("</think>", "");

  const callStream = async () => {
    let data = "";
    controller = new AbortController();

    try {
      const response = await fetch(API_URL, {
        signal: controller.signal,
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ query }),
      });

      if (!response.body) {
        console.error("ReadableStream not supported");
        return;
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();

      while (true) {
        const { value, done } = await reader.read();
        if (done) {
          console.info("Done reading.");
          console.info({ value });
          break;
        }

        const chunk = decoder.decode(value, { stream: true });
        if (!chunk) continue;

        data += chunk;

        setCurrentAnswer(removeThinkTag(data));
      }
    } catch (err) {
      console.error(err);
    }

    return removeThinkTag(data);
  };

  const callExtract = async () => {
    let data = "";
    controller = new AbortController();

    try {
      const response = await fetch(`${API_URL}extract`, {
        signal: controller.signal,
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ query }),
      });

      data = await response.text();
    } catch (err) {
      console.error(err);
    }

    return data;
  };

  const onSearchClick = async () => {
    const t0 = performance.now();

    setThinking(true);

    let data = "";

    if (extract) {
      data = await callExtract();
    } else {
      data = (await callStream()) ?? "";
    }

    // calculate benchmark
    const t1 = performance.now();
    const seconds = (t1 - t0) / 1_000;
    const secondsFormatted = `Took ${seconds.toFixed(2)} seconds`;

    setMessageElements((prev) => [
      ...prev,
      {
        answer: createMarkdownElement(data),
        time: secondsFormatted,
        question: createMarkdownElement(query),
      },
    ]);

    setThinking(false);
    setCurrentAnswer("");
  };

  const onAbortClick = () => {
    controller?.abort();
    setThinking(false);
  };

  const onInputChange = (val: string) => {
    setQuery(val);
  };

  const onCheckBoxChange = (checked: boolean) => {
    setExtract(checked);
  };

  const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.code === "Enter") {
      onSearchClick();
    }
  };

  return (
    <>
      <div className="container">
        <div id="messages" className="messages">
          {messageElements.map((elm, index) => (
            <React.Fragment key={index}>
              <div className="right-container">
                <div className="right">{elm.question}</div>
              </div>
              <div className="left">{elm.answer}</div>
              {elm.time && <div className="info">{elm.time}</div>}
            </React.Fragment>
          ))}

          {thinking && (
            <React.Fragment>
              <div className="right-container">
                <div className="right">{createMarkdownElement(query)}</div>
              </div>
              <div className="left">{createMarkdownElement(currentAnswer)}</div>
            </React.Fragment>
          )}
        </div>
        <Search
          inputRef={inputRef}
          thinking={thinking}
          onSearchClick={onSearchClick}
          onCheckBoxChange={onCheckBoxChange}
          onAbortClick={onAbortClick}
          onInputChange={onInputChange}
          onKeyDown={onKeyDown}
        />
      </div>
    </>
  );
}

export default App;
