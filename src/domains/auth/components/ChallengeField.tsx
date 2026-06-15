import { ReloadOutlined } from "@ant-design/icons";
import { Button, Input } from "antd";
import { useEffect, useState } from "react";
import {
  type AuthChallengeConfig,
  type AuthChallengeResponse,
} from "../api";
import { useImageChallengeQuery } from "../queries";
import styles from "./AuthScreen.module.css";

interface ChallengeFieldProps {
  config: AuthChallengeConfig;
  disabled?: boolean;
  resetKey: number;
  onChange(challenge?: AuthChallengeResponse): void;
  onError(message: string): void;
}

export function ChallengeField({
  config,
  disabled,
  resetKey,
  onChange,
  onError,
}: ChallengeFieldProps) {
  if (config.provider !== "image") {
    return (
      <div className={styles.remoteChallenge}>
        {config.sitekey ? "远程验证码已配置" : "远程验证码缺少站点公钥"}
      </div>
    );
  }

  return (
    <ImageChallengeField
      disabled={disabled}
      onChange={onChange}
      onError={onError}
      resetKey={resetKey}
    />
  );
}

interface ImageChallengeFieldProps {
  disabled?: boolean;
  resetKey: number;
  onChange(challenge?: AuthChallengeResponse): void;
  onError(message: string): void;
}

function ImageChallengeField({
  disabled,
  resetKey,
  onChange,
  onError,
}: ImageChallengeFieldProps) {
  const [answer, setAnswer] = useState("");
  const [refreshKey, setRefreshKey] = useState(0);
  const queryKey = resetKey + refreshKey;
  const challengeQuery = useImageChallengeQuery(queryKey, true);
  const challenge = challengeQuery.data;
  const loading = challengeQuery.isPending || challengeQuery.isFetching;

  useEffect(() => {
    setAnswer("");
    onChange(undefined);
    onError("");
  }, [onChange, onError, queryKey]);

  useEffect(() => {
    if (!challengeQuery.isError) {
      return;
    }
    onError(challengeQuery.error instanceof Error ? challengeQuery.error.message : "验证码生成失败");
  }, [challengeQuery.error, challengeQuery.isError, onError]);

  useEffect(() => {
    if (!challenge || answer.trim() === "") {
      onChange(undefined);
      return;
    }

    onChange({
      provider: "image",
      challengeId: challenge.challengeId,
      answer,
    });
  }, [answer, challenge, onChange]);

  return (
    <div className={styles.captchaRow}>
      <Input
        autoComplete="off"
        className={styles.captchaInput}
        disabled={disabled || loading}
        onChange={(event) => setAnswer(event.target.value)}
        placeholder="输入验证码"
        value={answer}
        variant="borderless"
      />
      <Button
        className={styles.captchaButton}
        disabled={disabled || loading}
        htmlType="button"
        onClick={() => setRefreshKey((value) => value + 1)}
      >
        {challenge ? <img alt="验证码" src={challenge.image} /> : <ReloadOutlined />}
      </Button>
    </div>
  );
}
