import { Image } from "expo-image";
import { useEffect, useRef, useState } from "react";
import {
  ScrollView,
  StyleSheet,
  useWindowDimensions,
  View,
  type NativeScrollEvent,
  type NativeSyntheticEvent,
} from "react-native";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

const BANNERS = [
  require("../../../assets/banners/banner_one.png"),
  require("../../../assets/banners/banner_two.png"),
  require("../../../assets/banners/banner_three.png"),
  require("../../../assets/banners/banner_four.png"),
];

const AUTO_ADVANCE_MS = 3500;
const HORIZONTAL_SCREEN_PADDING = layout.screenPaddingHorizontal * 2;

export function DashboardBannerCarousel() {
  const { width } = useWindowDimensions();
  const carouselRef = useRef<ScrollView>(null);
  const [activeSlide, setActiveSlide] = useState(0);
  const bannerWidth = Math.max(width - HORIZONTAL_SCREEN_PADDING, 0);

  useEffect(() => {
    if (bannerWidth <= 0) {
      return;
    }

    const timer = setInterval(() => {
      setActiveSlide((current) => {
        const next = (current + 1) % BANNERS.length;
        carouselRef.current?.scrollTo({
          x: next * bannerWidth,
          animated: true,
        });
        return next;
      });
    }, AUTO_ADVANCE_MS);

    return () => clearInterval(timer);
  }, [bannerWidth]);

  const handleMomentumScrollEnd = (
    event: NativeSyntheticEvent<NativeScrollEvent>,
  ) => {
    if (bannerWidth <= 0) {
      return;
    }

    setActiveSlide(Math.round(event.nativeEvent.contentOffset.x / bannerWidth));
  };

  return (
    <View style={styles.container}>
      <ScrollView
        ref={carouselRef}
        bounces={false}
        horizontal
        pagingEnabled
        showsHorizontalScrollIndicator={false}
        onMomentumScrollEnd={handleMomentumScrollEnd}
      >
        {BANNERS.map((source, index) => (
          <Image
            key={index}
            contentFit="contain"
            contentPosition="center"
            source={source}
            style={[styles.image, { width: bannerWidth }]}
          />
        ))}
      </ScrollView>

      <View style={styles.paginationWrap}>
        <View style={styles.pagination}>
          {BANNERS.map((_, index) => (
            <View
              key={index}
              style={[
                styles.dot,
                activeSlide === index ? styles.dotActive : styles.dotInactive,
              ]}
            />
          ))}
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: "#fff7f1",
    borderRadius: radius["2xl"],
    height: 132,
    overflow: "hidden",
    position: "relative",
    width: "100%",
    ...shadows.card,
  },
  dot: {
    borderRadius: radius.sm,
    height: spacing.sm,
  },
  dotActive: {
    backgroundColor: colors.text.inverse,
    width: spacing["2xl"],
  },
  dotInactive: {
    backgroundColor: "rgba(255, 255, 255, 0.55)",
    width: spacing.sm,
  },
  image: {
    height: "100%",
  },
  pagination: {
    backgroundColor: "rgba(0, 0, 0, 0.35)",
    borderRadius: radius["2xl"],
    flexDirection: "row",
    gap: spacing.sm,
    justifyContent: "center",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  paginationWrap: {
    alignItems: "center",
    bottom: spacing.sm,
    left: 0,
    position: "absolute",
    right: 0,
  },
});
